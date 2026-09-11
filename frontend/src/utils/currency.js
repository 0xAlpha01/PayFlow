const formatter = new Intl.NumberFormat("en-NG", {
  style: "currency",
  currency: "NGN",
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
});

export function formatCurrency(value) {
  return formatter.format(Number(value || 0));
}

export function formatCurrencyCompact(value) {
  const n = Number(value || 0);
  if (Math.abs(n) >= 1000) return `₦${(n / 1000).toLocaleString("en-NG", { maximumFractionDigits: 1 })}k`;
  return formatCurrency(n);
}
