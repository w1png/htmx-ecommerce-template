console.log("utils loaded");

function ClearFormOnSubmit(event, form) {
  if (!event.detail.successful || event.detail.xhr.status != 200) return;

  form.reset();
}

document.addEventListener("keydown", function(event) {
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

function OnChangeImage(input) {
  const img = document.getElementById("previewImage");
  img.src = URL.createObjectURL(input.files[0]);
  img.style.display = "block";
}

function OnChangeImages(input) {
  const img = document.getElementById("imagesPreviewList");
  img.innerHTML = "";

  if (input.files.length == 0) {
    img.classList.add("hidden");
    return;
  }

  for (const file of input.files) {
    const el = document.createElement("img");
    el.src = URL.createObjectURL(file);
    el.classList.add("admin-image-preview-list-item");
    img.appendChild(el);
  }

  img.classList.remove("hidden");
}

function ResetImageAndFormOnSubmit(event, form) {
  if (!event.detail.successful || event.detail.xhr.status != 200) return;

  form.reset();

  const img = document.getElementById("previewImage");
  img.src = "";
  img.style.display = "none";
}

function ResetImagesAndFormOnSubmit(event, form) {
  if (!event.detail.successful || event.detail.xhr.status != 200) return;

  form.reset();

  const images = document.getElementById("imagesPreviewList");
  images.innerHTML = "";

  images.classList.add("hidden");
}

htmx.on("htmx:afterSwap", function(evt) {
  if (evt.detail.target.id === "cart_products") {
    const checkout_products = document.getElementById("checkout_products");
    if (checkout_products) {
      checkout_products.innerHTML = evt.detail.target.innerHTML;
    }
  }
});

htmx.logAll();
