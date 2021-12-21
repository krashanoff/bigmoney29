/**
 * Check out the results from some upload.
 */

import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import Table from "react-bootstrap/Table";
import Alert from "react-bootstrap/Alert";
import Card from "react-bootstrap/Card";
import Stack from "react-bootstrap/Stack";

const Results = () => {
  const params = useParams();
  let sock = useRef(null);
  const [showMsg, setShowMsg] = useState(true);
  const [results, setResults] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    let websock = new WebSocket(`ws://localhost:8081/results/${params.id}`);
    websock.onopen = () => {
      console.info("Established connection to backend.");
      websock.send(params.id);
    };
    websock.onerror = (e) => {
      console.error(e);
      setError("Lost connection to the backend!");
    };
    websock.onmessage = (msg) => {
      setResults(JSON.parse(msg.data));
    };
    sock.current = websock;
    return () => websock.close();
  }, [params.id]);

  return (
    <Stack direction="vertical" gap={3}>
      {error ? (
        <section>
          <Alert variant="danger">{error}</Alert>
        </section>
      ) : (
        (showMsg || results.length === 0) && (
          <Alert
            variant="primary"
            dismissible
            onClick={() => setShowMsg(false)}
          >
            Your results weren't cached, so we're generating them at our
            earliest convenience.
          </Alert>
        )
      )}
      <Card body>
        <Card.Title>Results for run {params.id}</Card.Title>
        <Card.Subtitle className="text-muted">TODO: timestamp</Card.Subtitle>
      </Card>
      <Table striped bordered hover responsive>
        <thead>
          <tr>
            <th>#</th>
            <th>OK?</th>
            <th>Name</th>
            <th>Error Message</th>
          </tr>
        </thead>
        <tbody>
          {results &&
            Array.from(results).map((value, idx) => (
              <tr key={idx}>
                <td>{idx + 1}</td>
                <td>{value.fail ? "🟢" : "🔴"}</td>
                <td>TODO: name</td>
                <td>TODO: error message</td>
              </tr>
            ))}
        </tbody>
      </Table>
    </Stack>
  );
};

export default Results;
