import http from "k6/http";

export let options = {
  vus: 5,
  duration: "10s",
};

export default function () {
  http.get("http://host.docker.internal:8000");
  // request from inside container >> outside container
}

// k6 run script.js
// k6 run script.js -u2 -d10s
