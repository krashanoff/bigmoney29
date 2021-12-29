// Utilities for the API.

const USERLAND = "http://localhost:8081/largecurrency";

const login = (formEvent, onSuccess, onError) =>
  fetch(`${USERLAND}/login`, {
    method: "post",
    mode: "cors",
    body: new FormData(formEvent.target),
  })
    .then((r) => {
      if (r.status >= 200 && r.status < 300) return r.json();
      throw r.status;
    })
    .then((r) => onSuccess(r.token))
    .catch(onError);

const getAssignments = (jwt, onSuccess, onError) =>
  fetch(`${USERLAND}/assignments`, {
    method: "get",
    headers: {
      Authorization: `Bearer: ${jwt}`,
    },
  })
    .then((r) => {
      if (r.status < 300 && r.status >= 200) return r.json();
      throw r.status;
    })
    .then(onSuccess)
    .catch(onError);

const api = { login, getAssignments };
export default api;
