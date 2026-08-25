import { useRef } from "react";
export default function OTPInput({ value, onChange }) {
  const refs = useRef([]);
  const digits = value.padEnd(6, " ").slice(0, 6).split("");
  function update(index, digit) {
    const next = digits
      .map((item, i) => (i === index ? digit : item))
      .join("")
      .replaceAll(" ", "");
    onChange(next);
    if (digit && index < 5) refs.current[index + 1]?.focus();
  }
  function handleKey(index, event) {
    if (event.key === "Backspace" && !digits[index] && index > 0)
      refs.current[index - 1]?.focus();
  }
  function paste(event) {
    event.preventDefault();
    onChange(
      event.clipboardData.getData("text").replace(/\D/g, "").slice(0, 6),
    );
    refs.current[
      Math.min(event.clipboardData.getData("text").length, 5)
    ]?.focus();
  }
  return (
    <div className="otp-inputs" onPaste={paste}>
      {digits.map((digit, index) => (
        <input
          key={index}
          ref={(element) => {
            refs.current[index] = element;
          }}
          inputMode="numeric"
          maxLength={1}
          value={digit.trim()}
          onChange={(event) =>
            update(index, event.target.value.replace(/\D/g, "").slice(-1))
          }
          onKeyDown={(event) => handleKey(index, event)}
          aria-label={`Verification digit ${index + 1}`}
        />
      ))}
    </div>
  );
}
