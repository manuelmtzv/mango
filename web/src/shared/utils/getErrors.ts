import { isAxiosError } from "axios";
import { capitalizeFirstLetter } from "./capitalizeFirstLetter";

export function getErrors(error: unknown): string[] {
  if (isAxiosError(error)) {
    error = error.response?.data.message;

    if (typeof error === "string") {
      return [capitalizeFirstLetter(error)];
    }

    if (Array.isArray(error)) {
      return error.map((err) => capitalizeFirstLetter(err));
    }
  }

  if (error instanceof Error) {
    return [error.message];
  }

  if (typeof error === "string") {
    return [error];
  }

  return ["An unexpected error occurred. Please try again."];
}
