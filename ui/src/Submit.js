/**
 * Submit a file to the grading server, then redirect oneself to their
 * Results page.
 */

import { useNavigate } from "react-router-dom";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import Alert from "react-bootstrap/Alert";
import Stack from "react-bootstrap/Stack";

import "./Submit.css";
import { useState } from "react";

export default function Submit() {
  const navigate = useNavigate();
  const [errors, setErrors] = useState({
    uploadError: "",
    validationError: "",
  });

  const submit = (e) => {
    e.preventDefault();
    fetch("http://localhost:8081/upload", {
      method: "post",
      mode: "cors",
      body: new FormData(e.target),
    })
      .then((r) => {
        if (r.status >= 200 && r.status < 300) return r.text();
        throw r.status;
      })
      .then((resp) => navigate(`/${resp}`))
      .catch((e) =>
        setErrors((oldErr) => {
          oldErr.validationError = `Failed uploading! Server responded with: ${String(
            e
          ).replace("TypeError: ", "")}`;
          return oldErr;
        })
      );
  };

  return (
    <Stack direction="vertical" gap={3}>
      {errors.uploadError && (
        <Alert variant={"danger"}>
          Failed to send submission: {errors.uploadError}
        </Alert>
      )}
      <h1>Upload your file</h1>
      <p>
        This will upload your project to my server where it is automatically
        tested.
      </p>
      <Form onSubmit={submit}>
        <Form.Group className="mb-3" controlId="formEmail">
          <Form.Label>Name (Optional)</Form.Label>
          <Form.Control type="text" name="name" placeholder="Josie Bruin" />
        </Form.Group>
        <Form.Group className="mb-3" controlId="formAssignmentNumber">
          <Form.Label>Assignment</Form.Label>
          <Form.Select name="assignment" aria-label="Assignment chooser">
            <option>v0.1.0</option>
            <option value="1">v0.2.0</option>
            <option value="2">v0.3.0</option>
          </Form.Select>
        </Form.Group>
        <Form.Group className="mb-3" controlId="formFile">
          <Form.Label>
            <Form.Label>File (.tar.gz)</Form.Label>
            <Form.Control
              type="file"
              name="file"
              accept="application/gzip"
              required
            />
          </Form.Label>
        </Form.Group>
        <Button variant="primary" type="submit">
          Submit
        </Button>
      </Form>
    </Stack>
  );
}
