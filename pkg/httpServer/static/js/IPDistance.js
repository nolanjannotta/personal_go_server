const element = document.getElementById("distance")

async function getDistance() {
  const laLat = 34.052235
  const laLon = -118.243683

  const clientIPResp = await fetch("https://api.ipify.org?format=json")
  const clientIPJson = await clientIPResp.json()
  const clientLatLonResp = await fetch("http://ip-api.com/json/" + clientIPJson.ip)
  const clientLatLonJson = await clientLatLonResp.json()


  let lat1 = laLat * Math.PI / 180;
  let lon1 = laLon * Math.PI / 180;
  let lat2 = clientLatLonJson.lat * Math.PI / 180;
  let lon2 = clientLatLonJson.lon * Math.PI / 180;

  const km = Math.acos(Math.sin(lat1) * Math.sin(lat2) + Math.cos(lat1) * Math.cos(lat2) * Math.cos(lon2 - lon1)) * 6371;
  // const miles = Math.floor(km * 0.621371)
  element.textContent = Math.floor(km) + " km away"


}

document.onload = getDistance()