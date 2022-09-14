import { useNavigate } from "react-router-dom";
import { useState } from "react";
import { useCookies } from "react-cookie";

import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import Alert from "react-bootstrap/Alert";
import Stack from "react-bootstrap/Stack";
import Container from "react-bootstrap/Container";

import api from "./api";

const Login = ({ className = "[YOUR CLASS]" }) => {
  const navigate = useNavigate();
  const [cookies, setCookies] = useCookies(["jwt"]);
  const [error, setError] = useState("");

  if (cookies.jwt) navigate("/me");

  const submit = (e) => {
    e.preventDefault();
    api.login(
      e,
      (token) => {
        setCookies("jwt", token, "/");
        navigate("/me");
      },
      (e) =>
        setError(
          `Failed uploading! Server responded with: ${String(e).replace(
            "TypeError: ",
            ""
          )}`
        )
    );
  };

  return (
    <Container fluid="md">
      <Stack direction="vertical" gap={3}>
        {error && <Alert variant={"danger"}>Failed to login: {error}</Alert>}
        <h1>Log into {className}</h1>
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
    </Container>
  );
};

export default Login;
