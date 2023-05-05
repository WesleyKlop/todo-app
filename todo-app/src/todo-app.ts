import { LitElement, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import litLogo from "./assets/lit.svg";
import viteLogo from "/vite.svg";

/**
 * An example element.
 *
 * @slot - This element has a slot
 * @csspart button - The button
 */
@customElement("todo-app")
export class TodoApp extends LitElement {
  /**
   * The number of times the button has been clicked.
   */
  @property({ type: Number })
  count = 0;

  render() {
    return html`
      <button @click=${this._onClick} part="button">
        count is ${this.count}
      </button>
    `;
  }

  private _onClick() {
    this.count++;
  }

  static styles = css`
    :host {
      background-color: rebeccapurple;
    }
  `;
}

declare global {
  interface HTMLElementTagNameMap {
    "todo-app": TodoApp;
  }
}
