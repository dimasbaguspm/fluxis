/**
 * Utility function to combine classnames conditionally
 * Filters out falsy values and joins with spaces
 */
export const cx = (...classes: (string | undefined | null | false)[]): string => {
  return classes.filter(Boolean).join(" ");
};
