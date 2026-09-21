import type { Component } from "solid-js";

type UsernameProps = {
  value: string;
  onInput: (value: string) => void;
};

const Username: Component<UsernameProps> = (props) => {
  return (
    <label class="input validator">
      <svg class="h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
        <g stroke-linejoin="round" stroke-linecap="round" stroke-width="2.5" fill="none" stroke="currentColor">
          <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
          <circle cx="12" cy="7" r="4"></circle>
        </g>
      </svg>
      <input
        type="text"
        required
        placeholder="Username"
        pattern="[A-Za-z][A-Za-z0-9\-]*"
        minlength="3"
        maxlength="30"
        title="Only letters, numbers or dash"
        value={props.value}
        onInput={(e) => props.onInput(e.currentTarget.value)}
      />
    </label>
  );
};

export default Username;