import {
  getRouteName,
  getLocationName,
  getVehicleName,
  getManufacturerName,
  getClassName,
  getSessionsByDay,
  getCurrentDriver,
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

  Alpine.store("state", {
    currentDate: "2025-01-01",
  });

  Alpine.store("currentDriver", {
    name: "N/A",
  });

  Alpine.store("driverAttempts", {
    selectedDriver: null,
    driverAttempts: [],
  });

  Alpine.data("controller", () => ({
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
      Alpine.store("state").currentDate = newDate;
      Alpine.store("event").date = formatDate(newDate);
      fetchSessionsForDay(newDate);
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
        return;
      }

      const trackCounts = {};
      const carClassCounts = {};
      data.forEach((item) => {
        trackCounts[item.route_id] = (trackCounts[item.route_id] ?? 0) + 1;
        carClassCounts[item.vehicle_class_id] =
          (carClassCounts[item.vehicle_class_id] ?? 0) + 1;
      });

      const mainTrackId = Object.keys(trackCounts).reduce((a, b) =>
        trackCounts[a] > trackCounts[b] ? a : b
      );
      const mainCarClassId = Object.keys(carClassCounts).reduce((a, b) =>
        carClassCounts[a] > carClassCounts[b] ? a : b
      );

      const filteredData = data.filter(
        (item) =>
          item.route_id == mainTrackId &&
          item.vehicle_class_id == mainCarClassId
      );

      // Fetch event info
      const routeName = await getRouteName(mainTrackId);
      const locationName = await getLocationName(filteredData[0]?.location_id);
      const carName = await getVehicleName(filteredData[0]?.vehicle_id);
      const manufacturerName = await getManufacturerName(
        filteredData[0]?.vehicle_manufacturer_id
      );
      const className = await getClassName(mainCarClassId);

      // Update the event info
      Alpine.store("event", {
        class: className,
        car: carName,
        manufacturer: manufacturerName,
        date: formatDate(date),
        stage: routeName,
        location: locationName,
      });

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

  // Initialize the page with today's date
  const today = new Date().toISOString().split("T")[0];
  Alpine.store("state").currentDate = today;
  fetchCurrentDriver();
  fetchSessionsForDay(today);

  setInterval(() => {
    fetchCurrentDriver();
    fetchSessionsForDay(Alpine.store("state").currentDate);
  }, 5000);
});
