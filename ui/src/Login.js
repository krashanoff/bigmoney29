import { useNavigate } from "react-router-dom";
import { useState } from "react";
import { useCookies } from 'react-cookie';

import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import Alert from "react-bootstrap/Alert";
import Stack from "react-bootstrap/Stack";

const Login = () => {
  const navigate = useNavigate();
  const [cookies, setCookies] = useCookies(["jwt"]);
  const [error, setError] = useState("");

  if (cookies.jwt)
    navigate("/me");

  const submit = (e) => {
    e.preventDefault();
    fetch("http://localhost:8081/cash/login", {
      method: "post",
      mode: "cors",
      body: new FormData(e.target),
    })
      .then((r) => {
        if (r.status >= 200 && r.status < 300) return r.json();
        throw r.status;
      })
      .then((resp) => {
        setCookies("jwt", resp.token, "/");
        navigate("/me");
      })
      .catch((e) =>
        setError(
          `Failed uploading! Server responded with: ${String(
            e
          ).replace("TypeError: ", "")}`
        )
      );
  };

  return (
    <Stack direction="vertical" gap={3}>
      {error && (
        <Alert variant={"danger"}>
          Failed to login: {error}
        </Alert>
      )}
      <h1>Your class name</h1>
      <p>Log in to the [YOUR CLASS] grading system.</p>
      <Form onSubmit={submit}>
        <Form.Group className="mb-3" controlId="formUsername">
          <Form.Label>Username</Form.Label>
          <Form.Control type="text" name="username" placeholder="Josie Bruin" />
        </Form.Group>
        <Form.Group className="mb-3" controlId="formPassword">
          <Form.Label>Password</Form.Label>
          <Form.Control type="password" name="password" />
        </Form.Group>
        <Button variant="primary" type="submit">
          Login
        </Button>
      </Form>
    </Stack>
  );
};

export default Login;
