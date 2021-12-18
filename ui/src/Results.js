import { useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import Table from 'react-bootstrap/Table';
import Alert from 'react-bootstrap/Alert';

const Results = () => {
  const params = useParams();
  let sock = useRef(null);
  const [showMsg, setShowMsg] = useState(true);
  const [results, setResults] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    let websock = new WebSocket(`ws://localhost:8080/results/${params.id}`);
    websock.onopen = () => {
      console.info("Established connection to backend.")
      websock.send(params.id);
    };
    websock.onerror = () => {
      setError("Lost connection to the backend!");
    };
    websock.onmessage = (msg) => {
      setResults(msg.data);
    };
    sock.current = websock;
    return () => websock.close();
  }, [params.id]);

  return (
    <>
      {error ?
        <section>
          {error}
        </section>
        : ((showMsg || results.length === 0) &&
          <Alert variant='primary' dismissible onClick={() => setShowMsg(false)}>
            Your results weren't cached, so we're generating them at our earliest convenience.
          </Alert>
        )}
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
          {Array.from(results).map((value, idx) => (
            <tr key={idx}>
              <td>{idx + 1}</td>
              <td>{String(value)}</td>
            </tr>
          ))}
        </tbody>
      </Table>
    </>
  );
}

export default Results;