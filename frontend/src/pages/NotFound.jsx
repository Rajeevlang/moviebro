import { Link } from "react-router-dom";
import { EmptyState } from "../components/ui";

export default function NotFound() {
  return (
    <div className="page">
      <EmptyState icon="🛸" title="404 — Page not found" hint="That page drifted off into deep space." />
      <p className="center">
        <Link to="/" className="btn btn-primary">
          Back to movies
        </Link>
      </p>
    </div>
  );
}
