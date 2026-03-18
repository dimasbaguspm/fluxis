import dayjs, { type Dayjs } from "dayjs";

export enum FormatDate {
  ShortDate = "MM/DD/YYYY",
  LongDate = "MMMM D, YYYY",
  ShortTime = "HH:mm",
  LongTime = "HH:mm:ss",
  DateTime = "MM/DD/YYYY HH:mm",
  DateTimeSeconds = "MM/DD/YYYY HH:mm:ss",
  ISO = "YYYY-MM-DDTHH:mm:ss[Z]",
  DateInput = "YYYY-MM-DD",
}

export function dateFormat(
  date: string | Date | Dayjs | null | undefined,
  format: FormatDate,
): string {
  if (!date) {
    return "";
  }

  return dayjs(date).format(format);
}
