import type { ErrorResponse } from "~/types/server-errors";

type MapServerErrorsArg = {
  errors: ErrorResponse["errors"];
  messages: Record<string, string>;
  fields?: Record<string, Record<string, string>>;
};

type MapServerErrorsReturn<T extends MapServerErrorsArg> = {
  messages: string[];
  fields: Record<keyof T["fields"], string[]>;
};

export function mapServerErrors<T extends MapServerErrorsArg>({
  errors,
  messages,
  fields,
}: T): MapServerErrorsReturn<T> {
  const messageErrors: string[] = [];

  const fieldsArray = Object.entries(fields ?? {});
  const fieldErrors = fieldsArray.reduce<MapServerErrorsReturn<T>["fields"]>(
    (acc, [field]) => {
      acc[field as keyof typeof acc] = [];
      return acc;
    },
    {} as MapServerErrorsReturn<T>["fields"],
  );

  for (const error of errors) {
    const genericRemap = messages[error.message];

    if (genericRemap) {
      messageErrors.push(genericRemap);
      break;
    }

    const field = fieldsArray.find(([field]) => field === error.field);

    if (field) {
      const [fieldName, messageRemaps] = field;
      fieldErrors[fieldName as keyof typeof fieldErrors].push(
        messageRemaps[error.message] ?? error.message,
      );
    }
  }

  return {
    messages: messageErrors,
    fields: fieldErrors as MapServerErrorsReturn<T>["fields"],
  };
}
