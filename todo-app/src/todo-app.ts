import { LitElement, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { map } from "lit/directives/map.js";

type Todo = {
  content: string;
  id: string;
};

@customElement("todo-app")
export class TodoApp extends LitElement {
  @property()
  todos: Todo[] = [];

  render() {
    return html`
      <ul>
        ${map(this.todos, (todo: Todo) => html`<li>${todo.content}</li>`)}
      </ul>
      <form method="POST" @submit=${this.submitTodo}>
        <input type="text" name="content" />
        <button type="submit">maak todo</button>
      </form>
    `;
  }

  submitTodo(evt: SubmitEvent) {
    evt.preventDefault();
    const formData = new FormData(evt.currentTarget as HTMLFormElement);
    const content = formData.get("content") || "";
    if (typeof content !== "string") {
      throw new Error("How is this not a string");
    }
    this.createTodo(content);
  }

  connectedCallback(): void {
    super.connectedCallback();
    this.fetchTodos();
  }

  async fetchTodos() {
    const response = await fetch("/api/todos", {
      headers: {
        accept: "application/json",
      },
    });
    this.todos = await response.json();
  }

  async createTodo(content: string) {
    const response = await fetch("/api/todos", {
      method: "POST",
      headers: {
        accept: "application/json",
        "content-type": "application/json",
      },
      body: JSON.stringify({ content }),
    });
    if (response.status !== 201) {
      console.error("Failed to create todo smh", response);
    }
    await this.fetchTodos();
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
