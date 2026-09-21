import { createSignal } from "solid-js";
import Password from "./Password"
import Username from "./Username"
import { usersApi } from "../../api/auth";
import { action, useAction, useSubmission } from "@solidjs/router";


export const loginAction = action(async (username: string, password: string) => {
  return usersApi.login({ username, password });
});

type LoginProps = {
  onGoToSignup: () => void;
  onSuccess?: (data: unknown) => void
};

const LogIn = (props: LoginProps) => {
  const [username, setUsername] = createSignal("");
  const [password, setPassword] = createSignal("");

  const login = useAction(loginAction);
  const submission = useSubmission(loginAction);

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    const user = await login(username(), password());
    if (user) props.onSuccess?.(user);
  };

  return (
    <div class="flex align-items justify-center h-screen">
      <div class="aura aura-dual m-auto">
        <div class="card bg-base-100">
          <div class="card-body">
            <div class="flex justify-center">
              <span class="text-base">Welcome to <b>Solanum Streaming</b></span>
            </div>

            <form onSubmit={handleSubmit}>
              <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
                <legend class="fieldset-legend text-primary">Log In</legend>

                <Username value={username()} onInput={setUsername} />
                <Password value={password()} onInput={setPassword} />

                <div class="tooltip tooltip-open tooltip-top tooltip-end mt-4 mb-2" data-tip={submission.error ? submission.error.message : ""}>
                  <button class="btn btn-soft btn-primary w-full" type="submit" disabled={submission.pending}>
                    {submission.pending ? "Logging in..." : "Log In"}
                  </button>
                </div>
              </fieldset>
            </form>
          </div>
          <div class="flex items-center justify-center gap-1">
            <p class="text-xs">Need an account?</p>
            <button class="btn btn-xs btn-link" onClick={props.onGoToSignup}>
              Sign Up
            </button>
          </div>
          <br/>
        </div>
      </div>
    </div>
  )
}

export default LogIn;