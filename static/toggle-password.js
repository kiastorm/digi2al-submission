customElements.define(
  "toggle-password",
  class extends HTMLElement {
    /**
     * Instantiate the Web Component
     */
    constructor() {
      super();

      this.passwords = this.querySelectorAll('[type="password"]');
      this.trigger = this.querySelector("[toggle]");
      if (!this.trigger) return;
      this.visible = this.hasAttribute("visible");

      this.init();
    }

    /**
     * Handle events
     * @param  {Event} event The event object
     */
    handleEvent(event) {
      this.toggle();
    }

    /**
     * Show hidden elements and add ARIA
     */
    init() {
      const hidden = this.trigger.closest("[hidden]");

      if (hidden) {
        hidden.removeAttribute("hidden");
      }

      this.trigger.setAttribute("aria-pressed", this.visible);
      this.trigger.setAttribute("type", "button");

      if (this.visible) {
        this.show();
      }

      this.trigger.addEventListener("click", this);
    }

    /**
     * Show passwords
     */
    show() {
      for (let pw of this.passwords) {
        pw.type = "text";
      }

      this.trigger.setAttribute("aria-pressed", true);
    }

    /**
     * Hide password visibility
     */
    hide() {
      for (let pw of this.passwords) {
        pw.type = "password";
      }

      this.trigger.setAttribute("aria-pressed", false);
    }

    /**
     * Toggle password visibility on or off
     */
    toggle() {
      const show = this.trigger.getAttribute("aria-pressed") === "false";

      if (show) {
        this.show();
      } else {
        this.hide();
      }
    }
  }
);
