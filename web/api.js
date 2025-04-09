export async function postQuery(queryPayload) {
  const response = await fetch("/api/query", {
    method: "POST",
    headers: {
      "Content-Type": "plain/text; charset=utf-8",
    },
    body: queryPayload,
  });
  if (!response.ok) {
    throw new Error("Failed to postQuery", queryPayload);
  }
  const data = await response.json();
  return data;
}

// Generic function to fetch a single name from a table
async function getNameById(
  tableName,
  id,
  idColumn = "id",
  defaultValue = "Unknown"
) {
  try {
    const queryPayload = `SELECT name FROM ${tableName} WHERE ${idColumn} = ${id}`;
    const data = await postQuery(queryPayload);
    return data[0]?.name || defaultValue;
  } catch (error) {
    console.error(`Error fetching name from ${tableName}:`, error);
    return defaultValue;
  }
}

// Specific functions using the generic getNameById
export async function getRouteName(routeId) {
  return getNameById("routes", routeId, "id", "Unknown Route");
}

export async function getLocationName(locationId) {
  return getNameById("locations", locationId, "id", "Unknown Location");
}

export async function getVehicleName(vehicleId) {
  return getNameById("vehicles", vehicleId, "id", "Unknown Vehicle");
}

export async function getManufacturerName(manufacturerId) {
  return getNameById(
    "vehicle_manufacturers",
    manufacturerId,
    "id",
    "Unknown Manufacturer"
  );
}

export async function getClassName(classId) {
  return getNameById("vehicle_classes", classId, "id", "Unknown Class");
}

// Fetch sessions by day
export async function getSessionsByDay(date) {
  try {
    const startOfDay = `${date} 00:00:00`;
    const endOfDay = `${date} 23:59:59`;
    const queryPayload = `SELECT sessions.user_id, users.name AS user_name, 
                    sessions.stage_result_time AS time,
                    sessions.route_id, sessions.location_id, 
                    sessions.vehicle_id, sessions.vehicle_manufacturer_id, 
                    sessions.vehicle_class_id,
                    sessions.stage_result_status
            FROM sessions 
            JOIN users ON sessions.user_id = users.id 
            WHERE sessions.started_at BETWEEN '${startOfDay}' AND '${endOfDay}'`;

    const data = await postQuery(queryPayload);
    return data;
  } catch (error) {
    console.error("Error fetching sessions:", error);
    return [];
  }
}

// Fetch the current driver
export async function getCurrentDriver() {
  try {
    const queryPayload = `SELECT users.name AS user_name
              FROM user_logins
              JOIN users ON user_logins.user_id = users.id
              ORDER BY user_logins.timestamp DESC
              LIMIT 1;`;
    const data = await postQuery(queryPayload);
    return data[0]?.user_name || "N/A";
  } catch (error) {
    console.error("Error fetching current driver:", error);
    return "N/A";
  }
}
