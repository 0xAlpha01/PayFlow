import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { CheckCircle2, Clock, XCircle } from "lucide-react";
import { verifyPayment } from "../services/paymentService";

const money = (value, currency = "NGN") =>
  new Intl.NumberFormat("en-NG", { style: "currency", currency }).format(Number(value || 0));

export default function PaymentCallback() {
  const [params] = useSearchParams();
  const reference = params.get("reference") || params.get("trxref") || "";
  const [status, setStatus] = useState("checking");
  const [payment, setPayment] = useState(null);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!reference) {
      setStatus("error");
      setError("No payment reference was returned by Paystack.");
      return;
    }
    let cancelled = false;
    verifyPayment(reference)
      .then(({ data }) => {
        if (cancelled) return;
        setPayment(data.data);
        setStatus(data.data.status === "successful" ? "success" : data.data.status === "failed" ? "failed" : "pending");
      })
      .catch((requestError) => {
        if (cancelled) return;
        setError(requestError.response?.data?.message || "Unable to verify this payment.");
        setStatus("error");
      });
    return () => {
      cancelled = true;
    };
  }, [reference]);

  return (
    <>
      <div className="page-heading">
        <div>
          <h2>Payment result</h2>
          <p>We confirm every payment directly with Paystack before updating your records.</p>
        </div>
      </div>
      <div className="panel" style={{ padding: 40, textAlign: "center" }}>
        {status === "checking" && <div className="data-loading">Verifying your payment...</div>}

        {status === "success" && (
          <>
            <CheckCircle2 size={40} color="#16a66a" />
            <h3 style={{ margin: "16px 0 6px" }}>Payment successful</h3>
            <div className="auth-alert success" style={{ display: "inline-block", margin: "0 0 18px" }}>
              {money(payment.amount, payment.currency)} was received and your transaction has been recorded.
            </div>
          </>
        )}

        {status === "pending" && (
          <>
            <Clock size={40} color="#b98900" />
            <h3 style={{ margin: "16px 0 6px" }}>Payment pending</h3>
            <div className="auth-alert" style={{ display: "inline-block", margin: "0 0 18px", color: "#8a6d00", background: "#fff7e0", border: "1px solid #f2e2a8" }}>
              Paystack hasn't confirmed this payment yet. It will update automatically once processed.
            </div>
          </>
        )}

        {status === "failed" && (
          <>
            <XCircle size={40} color="#c0392b" />
            <h3 style={{ margin: "16px 0 6px" }}>Payment failed</h3>
            <div className="auth-alert error" style={{ display: "inline-block", margin: "0 0 18px" }}>
              This payment was not completed. You can try again from the Payments page.
            </div>
          </>
        )}

        {status === "error" && (
          <>
            <XCircle size={40} color="#c0392b" />
            <h3 style={{ margin: "16px 0 6px" }}>Unable to verify payment</h3>
            <div className="auth-alert error" style={{ display: "inline-block", margin: "0 0 18px" }}>{error}</div>
          </>
        )}

        {payment && (
          <p style={{ color: "var(--muted)", fontSize: 12, marginBottom: 24 }}>
            Reference: <strong>{payment.paystackReference}</strong>
          </p>
        )}

        {status !== "checking" && (
          <div style={{ display: "flex", gap: 10, justifyContent: "center" }}>
            <Link className="button" to="/payments">View payments</Link>
            <Link className="button button-secondary" to="/transactions">View transactions</Link>
          </div>
        )}
      </div>
    </>
  );
}
