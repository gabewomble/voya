import { component$ } from "@builder.io/qwik";
import type { FieldElementProps, FieldStore } from "@modular-forms/qwik";

type TextInputProps = {
  id: string;
  name: string;
  label: string;
  placeholder?: string;
  value?: string;
  type?: "text" | "password" | "textarea";
  fieldProps: FieldElementProps<Record<string, string>, string>;
  field: FieldStore<any, any>;
};

export const TextInput = component$<TextInputProps>(
  ({ id, label, name, placeholder, type = "text", fieldProps, field }) => (
    <>
      <label for={id} class="mb-2 block text-sm font-bold text-base-content">
        {label}
      </label>

      {type === "textarea" ? (
        <textarea
          {...fieldProps}
          class={`input input-bordered w-full ${field.error ? "input-error" : ""}`}
          id={id}
          name={name}
          placeholder={placeholder}
          value={field.value}
        />
      ) : (
        <input
          {...fieldProps}
          class={`input input-bordered w-full ${field.error ? "input-error" : ""}`}
          type={type}
          id={id}
          name={name}
          placeholder={placeholder}
          value={field.value}
        />
      )}
      {field.error && <p class="my-1 text-error">{field.error}</p>}
    </>
  ),
);
