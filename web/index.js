import {
  getRouteName,
  getLocationName,
  getVehicleName,
  getManufacturerName,
  getClassName,
  getSessionsByDay,
  getCurrentDriver,
  getChampionshipStandings,
} from "./api.js";

document.addEventListener("alpine:init", () => {
  Alpine.store("laptimer", {
    laptimes: [],
    latestDriveIndex: -1,
  });

  Alpine.store("event", {
    class: "",
    car: "",
    manufacturer: "",
    date: "",
    stage: "",
    location: "",
  });

  Alpine.store("championship", {
    enabled: true,
    standings: [],
  });

  Alpine.store("state", {
    currentDate: "2025-01-01",
    selectedStageVehicle: null,
    stageVehicleOptions: [],
  });

  Alpine.store("currentDriver", {
    name: "N/A",
  });

  Alpine.store("driverAttempts", {
    selectedDriver: null,
    driverAttempts: [],
  });

  Alpine.data("controller", () => ({
    showDropdown: false,

    async showDriverAttempts(driverName) {
      Alpine.store("driverAttempts").selectedDriver = driverName;

      // Show loading state
      Alpine.store("driverAttempts").driverAttempts = await new Promise(
        async (resolve) => {
          const allSessions = await getSessionsByDay(
            Alpine.store("state").currentDate
          );

          const attempts = allSessions
            .filter((session) => session.user_name === driverName)
            .map((attempt) => ({
              time:
                attempt.stage_result_status === 1
                  ? formatTime(attempt.time)
                  : "DNF",
              status: attempt.stage_result_status === 1 ? "Completed" : "DNF",
              car: attempt.vehicle_name || "-",
              class: attempt.vehicle_class_id || "-",
            }));
          resolve(attempts);
        }
      );

      // Show the modal
      document.querySelector(".modal").style.display = "block";
    },

    closeModal() {
      Alpine.store("driverAttempts").selectedDriver = null;
      Alpine.store("driverAttempts").driverAttempts = [];
      document.querySelector(".modal").style.display = "none";
    },

    resetLeaderboard() {
      // Reset to the main leaderboard
      this.selectedDriver = null;
      fetchSessionsForDay(Alpine.store("state").currentDate);
    },

    changeDay(offset) {
      const currentDate = new Date(
        Alpine.store("state").currentDate.split(".").reverse().join("-")
      );
      currentDate.setDate(currentDate.getDate() + offset);
      const newDate = currentDate.toISOString().split("T")[0];

      // Reset the selected stage and vehicle combo
      Alpine.store("state").selectedStageVehicle = null;

      // Update the current date in the store
      Alpine.store("state").currentDate = newDate;
      Alpine.store("event").date = formatDate(newDate);

      // Fetch sessions for the new date
      fetchSessionsForDay(newDate);
    },

    async updateStageVehicle(selectedValue) {
      // Save the selected combo to the store
      Alpine.store("state").selectedStageVehicle = selectedValue;

      // Immediately refetch sessions for the current day
      await fetchSessionsForDay(Alpine.store("state").currentDate);
    },
  }));

  async function fetchCurrentDriver() {
    try {
      const currentDriverName = await getCurrentDriver();
      Alpine.store("currentDriver").name = currentDriverName;
    } catch (error) {
      console.error("Error fetching current driver:", error);
    }
  }
  async function fetchSessionsForDay(date) {
    try {
      const data = await getSessionsByDay(date);
      if (!data || data.length === 0) {
        Alpine.store("event", {
          class: "-",
          car: "-",
          manufacturer: "",
          date: Alpine.store("event").date,
          stage: "-",
          location: "-",
        });
        Alpine.store("laptimer").laptimes = [];
        Alpine.store("state").stageVehicleOptions = [];
        return;
      }

      // Group sessions by stage and vehicle class combinations
      const stageVehicleMap = {};
      data.forEach((item) => {
        const key = `${item.route_id}+${item.vehicle_class_id}`; // Use routeId + vehicleClassId
        if (!stageVehicleMap[key]) {
          stageVehicleMap[key] = [];
        }
        stageVehicleMap[key].push(item);
      });

      // Populate dropdown options
      const stageVehicleOptions = Object.keys(stageVehicleMap).map((key) => {
        const [routeId, vehicleClassId] = key.split("+");
        return { routeId, vehicleClassId }; // Include vehicleClassId in the options
      });
      Alpine.store("state").stageVehicleOptions = stageVehicleOptions;

      // Use the selected stage and vehicle class combo if available
      let selectedKey = Alpine.store("state").selectedStageVehicle;
      if (!selectedKey) {
        // Default to the most recent session if no combo is selected
        const mostRecentSession = data.reduce((latest, session) =>
          new Date(session.started_at) > new Date(latest.started_at)
            ? session
            : latest
        );
        selectedKey = `${mostRecentSession.route_id}+${mostRecentSession.vehicle_class_id}`;
        Alpine.store("state").selectedStageVehicle = selectedKey;
      }

      const [routeId, vehicleClassId] = selectedKey.split("+");

      // Fetch event info for the selected session
      const routeName = await getRouteName(routeId);
      const locationName = await getLocationName(
        stageVehicleMap[selectedKey][0]?.location_id
      );
      const className = await getClassName(vehicleClassId);

      // Update the event info
      Alpine.store("event", {
        class: className,
        car: "-", // No specific car, as it's based on vehicle class
        manufacturer: "-", // No specific manufacturer, as it's based on vehicle class
        date: formatDate(date),
        stage: routeName,
        location: locationName,
      });

      // Filter sessions based on the selected stage and vehicle class
      const filteredData = stageVehicleMap[selectedKey];

      // Group laps by user and calculate fastest lap and total attempts
      const userLaps = {};
      filteredData.forEach((item) => {
        const userId = item.user_id;
        const time = parseFloat(item.time);

        if (!userLaps[userId]) {
          userLaps[userId] = {
            name: item.user_name,
            fastestTime: item.stage_result_status === 1 ? time : Infinity,
            totalAttempts: 1,
            hasCompleted: item.stage_result_status === 1,
            carClass: item.vehicle_class_name, // Use vehicle class name
            vehicle: item.vehicle_name || "-", // Fallback to "-" if vehicle_name is empty
          };
        } else {
          if (item.stage_result_status === 1) {
            userLaps[userId].fastestTime = Math.min(
              userLaps[userId].fastestTime,
              time
            );
            userLaps[userId].hasCompleted = true;
          }
          userLaps[userId].totalAttempts += 1;
        }
      });

      // Fetch vehicle names if missing
      for (const userId in userLaps) {
        if (userLaps[userId].vehicle === "-") {
          const vehicleId = filteredData.find(
            (item) => item.user_id === userId
          )?.vehicle_id;

          if (vehicleId) {
            userLaps[userId].vehicle = await getVehicleName(vehicleId);
          }
        }
      }

      const laptimes = Object.values(userLaps)
        .filter((user) => user.hasCompleted || user.fastestTime === Infinity) // Include DNF only if no completed runs
        .map((user, index) => ({
          position: 0,
          name: user.name,
          time:
            user.fastestTime === Infinity
              ? "DNF"
              : formatTime(user.fastestTime),
          rawTime: user.fastestTime,
          diff: "",
          attempts: user.totalAttempts,
          carClass: user.carClass, // Include carClass in the laptime object
          vehicle: user.vehicle, // Include the specific car used
        }))
        .sort((a, b) => {
          if (a.rawTime === Infinity) return 1; // Place DNF at the bottom
          if (b.rawTime === Infinity) return -1;
          return a.rawTime - b.rawTime;
        });

      const fastestTime = laptimes[0]?.rawTime ?? 0;
      laptimes.forEach((laptime, index) => {
        laptime.position = index + 1;
        laptime.diff =
          laptime.rawTime > fastestTime && laptime.rawTime !== Infinity
            ? `+${formatTime(laptime.rawTime - fastestTime)}`
            : "";
      });
      Alpine.store("laptimer").laptimes = laptimes;
    } catch (error) {
      console.error("Error parsing sessions:", error);
      Alpine.store("event", {
        class: "-",
        car: "-",
        manufacturer: "",
        date: Alpine.store("event").date,
        stage: "-",
        location: "-",
      });
      Alpine.store("laptimer").laptimes = [];
      Alpine.store("state").stageVehicleOptions = [];
    }
  }

  function formatTime(seconds) {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    const ms = Math.round((seconds % 1) * 1000);
    return `${String(mins).padStart(2, "0")}:${String(secs).padStart(
      2,
      "0"
    )}.${String(ms).padStart(3, "0")}`;
  }

  function formatDate(date) {
    const [year, month, day] = date.split("-");
    return `${day}.${month}`;
  }

  // Extract the "championship" parameter from the URL
  const urlParams = new URLSearchParams(window.location.search);
  const championshipId = urlParams.get("championship");

  // Initialize the page with today's date
  const today = new Date().toISOString().split("T")[0];
  Alpine.store("state").currentDate = today;
  fetchCurrentDriver();
  fetchSessionsForDay(today);

  // Fetch championship standings if the "championship" parameter exists
  if (championshipId) {
    fetchChampionshipStandings(championshipId);
  }

  // Set up periodic updates
  setInterval(() => {
    fetchCurrentDriver();
    fetchSessionsForDay(Alpine.store("state").currentDate);

    // Fetch championship standings only if the "championship" parameter exists
    if (championshipId) {
      fetchChampionshipStandings(championshipId);
    }
  }, 5000);
});

document.addEventListener("DOMContentLoaded", () => {
  // Check for the "championship" parameter in the URL
  const urlParams = new URLSearchParams(window.location.search);
  const championshipId = urlParams.get("championship");

  if (championshipId) {
    // Show the championship-results section
    const championshipResults = document.getElementById("championship-results");
    if (championshipResults) {
      championshipResults.style.display = "block";
    }

    // Fetch championship standings using the parameter
    fetchChampionshipStandings(championshipId);
  }
});

async function fetchChampionshipStandings(id) {
  try {
    const standings = await getChampionshipStandings(id);
    Alpine.store("championship").standings = standings.map(
      (standing, index) => ({
        position: index + 1,
        name: standing.user_name,
        points: standing.points,
      })
    );
  } catch (error) {
    console.error("Error fetching championship standings:", error);
  }
}
