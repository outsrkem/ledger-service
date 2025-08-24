// Time format filter
export function formatTime(value, format = "YYYY-MM-DD HH:mm:ss") {
    // Handle invalid timestamps
    if (value === null || value === "" || value === undefined || value <= 0) {
        return "--";
    }
    // Handle second-level timestamps (if less than 1e12, treat as seconds and convert to milliseconds)
    if (value < 1e12) {
        value = value * 1000;
    }
    const date = new Date(value);
    // Get individual time components
    const year = date.getFullYear();
    const month = date.getMonth() + 1; // Months start from 0
    const day = date.getDate();
    const hours = date.getHours();
    const minutes = date.getMinutes();
    const seconds = date.getSeconds();
    // Zero-padding function
    const padZero = (num) => (num < 10 ? "0" + num : num);
    // Replace placeholders in format string
    return (
        format
            .replace("YYYY", year)
            .replace("MM", padZero(month))
            .replace("DD", padZero(day))
            .replace("HH", padZero(hours))
            .replace("mm", padZero(minutes))
            .replace("ss", padZero(seconds))
            // Support short formats
            .replace("YY", String(year).slice(2))
            .replace("M", month)
            .replace("D", day)
            .replace("H", hours)
            .replace("m", minutes)
            .replace("s", seconds)
    );
}
