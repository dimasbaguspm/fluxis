/**
 * Extract 2-character initials from a name
 * @param name - Full name or display name
 * @returns 2-character initials (uppercase)
 */
export const getInitials = (name: string | null | undefined): string => {
  if (!name) return "U";

  const trimmed = name.trim();
  if (trimmed.length === 0) return "U";

  const parts = trimmed.split(/\s+/);

  if (parts.length === 1) {
    return trimmed.substring(0, 2).toUpperCase();
  }

  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
};
