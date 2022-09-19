import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

import { useSelector, useDispatch } from "react-redux";
import { updateAssignments } from "./features/assignments";

import Card from "react-bootstrap/Card";
import Container from "react-bootstrap/Container";
import Placeholder from "react-bootstrap/Placeholder";
import ProgressBar from "react-bootstrap/ProgressBar";

const AssignmentCard = ({
  template = false,
  id,
  title,
  description,
  best_score,
  total_points,
}) => {
  const navigate = useNavigate();

  if (template) {
    return (
      <Card style={{ width: "18rem", margin: "2rem" }}>
        <Card.Body>
          <Placeholder as={Card.Title} animation="glow">
            <Placeholder xs={6} />
          </Placeholder>
          <Placeholder as={Card.Text} animation="glow">
            <Placeholder xs={7} /> <Placeholder xs={4} /> <Placeholder xs={4} />{" "}
            <Placeholder xs={6} /> <Placeholder xs={8} />
          </Placeholder>
        </Card.Body>
        <Card.Footer>
          <ProgressBar
            variant="primary"
            now={(best_score / total_points) * 100.0}
            label={`${best_score} / ${total_points} = ${
              best_score / total_points
            }`}
          />
        </Card.Footer>
      </Card>
    );
  }

  return (
    <Card
      style={{ width: "18rem", margin: "2rem" }}
      onClick={() => navigate(`/assignment/${id}`)}
    >
      <Card.Body>
        <Card.Title>{title}</Card.Title>
        <Card.Footer>Footer</Card.Footer>
      </Card.Body>
      <Card.Footer>
        <ProgressBar
          variant="primary"
          now={best_score}
          label={`${best_score} / ${total_points} = ${
            best_score / total_points
          }`}
        />
      </Card.Footer>
    </Card>
  );
};

const Landing = ({ className = "[YOUR CLASS]" }) => {
  const assn = useSelector((state) => state.assignments);
  const dispatch = useDispatch();

  useEffect(() => {
    const id = setInterval(() => {
      dispatch(updateAssignments());
    }, 3000);
    return () => clearTimeout(id);
  }, [assn]);

  return (
    <Container>
      <div
        style={{
          display: "flex",
          flexFlow: "row wrap",
        }}
      >
        {(assn.assignments || [1, 2, 3, 4, 5]).map((i, idx) => (
          <AssignmentCard
            template={!assn.assignments}
            key={idx}
            best_score={50}
            total_points={100}
            {...i}
          />
        ))}
      </div>
    </Container>
  );
};

export default Landing;
