import { useNavigate } from 'react-router-dom';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import Alert from 'react-bootstrap/Alert';
import Card from 'react-bootstrap/Card';

import './App.css';
import { useState } from 'react';

function App() {
  const navigate = useNavigate();
  const [error, setError] = useState("");

  const submit = (e) => {
    e.preventDefault();
    fetch("http://localhost:8080/upload", {
      method: "post",
      mode: 'cors',
      body: new FormData(e.target),
    })
      .then(r => {
        if (r.status >= 200 && r.status < 300)
          return r.text();
        throw r.status;
      })
      .then(resp => navigate(`/${resp}`))
      .catch(e => setError(`Failed uploading! Server responded with: ${String(e).replace("TypeError: ", "")}`));
  };

  return (
    <Card body>
      {error && <Alert variant={"danger"}>Failed to send submission: {error}</Alert>}
      <h1>Upload your file</h1>
      <Form onSubmit={submit}>
        <Form.Group className="mb-3" controlId='formEmail'>
          <Form.Label>Name</Form.Label>
          <Form.Control type="text" name="name" placeholder="Josie Bruin" />
        </Form.Group>
        <Form.Group className="mb-3" controlId="formAssignmentNumber">
          <Form.Label>Assignment</Form.Label>
          <Form.Select name="assignment" aria-label='Assignment chooser'>
            <option>v0.1.0</option>
            <option value="1">v0.2.0</option>
            <option value="2">v0.3.0</option>
          </Form.Select>
        </Form.Group>
        <Form.Group className="mb-3" controlId="formFile">
          <Form.Label>
            <Form.Label>File (.tar.gz)</Form.Label>
            <Form.Control type="file" name="file" />
          </Form.Label>
        </Form.Group>
        <Button variant="primary" type="submit">Submit</Button>
      </Form>
    </Card>
  );
}

export default App;
