function ClearFormOnSubmit(event, form) {
  if (!event.detail.successful || event.detail.xhr.status != 200) return;

  form.reset();
}

document.addEventListener("keydown", function (event) {
  if (event.key === "Escape") {
    const modals = document.querySelectorAll(".modal");
    const modal_overlays = document.querySelectorAll(".modal-overlay");

    modals.forEach((modal) => {
      modal.dataset.open = false;
    });

    modal_overlays.forEach((modal_overlay) => {
      modal_overlay.dataset.open = false;
    });
  }
});
