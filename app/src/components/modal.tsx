import { component$, Slot } from "@builder.io/qwik";

export type ModalProps = {
  id: string;
};

export const showModal = (id: string) => {
  const dialog = document.getElementById(id);
  if (dialog instanceof HTMLDialogElement) {
    dialog.showModal();
  }
};

export const hideModal = (id: string) => {
  const dialog = document.getElementById(id);
  if (dialog instanceof HTMLDialogElement) {
    dialog.close();
  }
};

export const ModalTitle = component$(() => {
  return (
    <h3 class="mb-4 text-xl font-bold">
      <Slot />
    </h3>
  );
});

export const ModalActions = component$(() => {
  return (
    <form method="dialog" class="flex gap-4">
      <Slot />
    </form>
  );
});

export const Modal = component$<ModalProps>(({ id }) => {
  return (
    <dialog id={id} class="modal">
      <div class="modal-box text-left text-lg">
        <Slot />
        <div class="modal-action">
          <Slot name="actions" />
        </div>
      </div>
    </dialog>
  );
});
