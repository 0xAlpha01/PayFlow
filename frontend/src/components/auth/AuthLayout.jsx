import { CreditCard, ShieldCheck, TrendingUp } from "lucide-react";
import { Link } from "react-router-dom";

export default function AuthLayout({ children, title, subtitle }) {
  return (
    <main className="auth-page">
      <section className="auth-brand-panel">
        <Link to="/login" className="auth-logo">
          <span className="brand-mark">P</span>
          <span>
            pay<span>flow</span>
          </span>
        </Link>
        <div className="auth-quote">
          <p>Payments that move your business forward.</p>
          <span>One clear view of your revenue, customers, and growth.</span>
        </div>
        <div className="auth-perks">
          <span>
            <TrendingUp size={16} /> Real-time insights
          </span>
          <span>
            <CreditCard size={16} /> Simple payments
          </span>
          <span>
            <ShieldCheck size={16} /> Built for trust
          </span>
        </div>
      </section>
      <section className="auth-form-panel">
        <div className="auth-form-wrap">
          <div className="auth-mobile-logo">
            <Link to="/login" className="auth-logo">
              <span className="brand-mark">P</span>
              <span>
                pay<span>flow</span>
              </span>
            </Link>
          </div>
          <div className="auth-heading">
            <h1>{title}</h1>
            <p>{subtitle}</p>
          </div>
          {children}
        </div>
      </section>
    </main>
  );
}
