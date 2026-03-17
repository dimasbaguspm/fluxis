import type { FC, ReactNode } from "react";

export interface WhenProps {
  condition: number | string | boolean | (number | boolean | string)[] | object | null | undefined;
  children: ReactNode | (() => ReactNode);
}

export const When: FC<WhenProps> = ({ condition, children }) => {
  const isArray = Array.isArray(condition);
  const isObject = typeof condition === "object" && condition !== null && !isArray;
  const isPrimitive = !isArray && !isObject;

  if (isArray) {
    const arr = condition as (number | boolean | string)[];
    if (arr.length === 0) return null;
    if (!arr.every(Boolean)) return null;
  }

  if (isObject) {
    const obj = condition as Record<string, unknown>;
    const keys = Object.keys(obj);
    if (keys.length === 0) return null;
    if (!Object.values(obj).every(Boolean)) return null;
  }

  if (isPrimitive && !condition) return null;

  const content = typeof children === "function" ? children() : children;
  return <>{content ?? null}</>;
};
