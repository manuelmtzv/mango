import { AxiosError } from "axios";

const defaultMessage = "Something went wrong! Please try again.";

export function getErrorMessage(
  error: unknown,
  fallbackMessage: string = defaultMessage,
): string {
  if (error instanceof AxiosError) {
    const { response } = error;

    const errMessage = response?.data?.message;

    if (typeof errMessage === "string") {
      return errMessage;
    } else if (Array.isArray(errMessage) && errMessage.length) {
      return errMessage[0];
    }
  } else if (error instanceof Error) {
    return error.message;
  } else if (typeof error === "string") {
    return error;
  }

  return fallbackMessage;
}
