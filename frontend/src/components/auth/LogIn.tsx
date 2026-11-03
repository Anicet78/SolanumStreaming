import { createSignal, type Component } from "solid-js";
import Password from "./Password"
import Username from "./Username"
import { usersApi } from "../../api/auth";
import { action, useAction, useSubmission } from "@solidjs/router";
import { useAuthContext } from "../../contexts/AuthContext";


export const loginAction = action(async (username: string, password: string) => {
  return usersApi.login({ username, password });
});

type LoginProps = {
  onGoToSignup: () => void;
  onSuccess?: (data: unknown) => void
};

const LogIn: Component<LoginProps> = (props) => {
  const auth = useAuthContext();

  const [username, setUsername] = createSignal("");
  const [password, setPassword] = createSignal("");

  const login = useAction(loginAction);
  const submission = useSubmission(loginAction);

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    const user = await login(username(), password());
    if (auth) auth.login(user);
  };

  return (
    <div class="auth-screen">
      <div class="auth-panel">
          <div class="auth-intro"><div class="eyebrow">Welcome back</div><strong>Solanum Streaming</strong><p>Your next favorite movie is waiting.</p></div>

            <form class="auth-form" onSubmit={handleSubmit}>
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
          <div class="auth-switch">
            <span>Need an account?</span>
            <button class="btn btn-xs btn-link" onClick={props.onGoToSignup}>
              Sign Up
            </button>
          </div>
      </div>
    </div>
  )
}

export default LogIn;
