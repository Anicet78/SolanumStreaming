import { createSignal, type Component } from "solid-js";
import Password from "./Password"
import PasswordConfirm from "./PasswordConfirm"
import Username from "./Username"
import { usersApi } from "../../api/auth";
import { action, useAction, useSubmission } from "@solidjs/router";
import { useAuthContext } from "../../contexts/AuthContext";

export const signupAction = action(async (username: string, password: string) => {
  return usersApi.create({ username, password });
});

type SignupProps = {
  onGoToLogin: () => void;
  onSuccess?: (data: unknown) => void
};

const SignUp: Component<SignupProps> = (props) => {
    const auth = useAuthContext();

  const [username, setUsername] = createSignal("");
  const [password, setPassword] = createSignal("");

  const login = useAction(signupAction);
  const submission = useSubmission(signupAction);

  const handleSubmit = async (e: SubmitEvent) => {
    e.preventDefault();
    const user = await login(username(), password());
    if (auth) auth.login(user);
  };

  return (
    <div class="auth-screen">
      <div class="auth-panel">
          <div class="auth-intro"><div class="eyebrow">Start your library</div><strong>Solanum Streaming</strong><p>Create an account and keep your watchlist close.</p></div>

            <form class="auth-form" onSubmit={handleSubmit}>
              <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
                <legend class="fieldset-legend text-primary">Sign Up</legend>

                <Username value={username()} onInput={setUsername} />
                <br/>
                <Password value={password()} onInput={setPassword} />
                <PasswordConfirm password={password}/>

                <div class="tooltip tooltip-open tooltip-top tooltip-end mt-4 mb-2" data-tip={submission.error ? submission.error.message : ""}>
                  <button class="btn btn-soft btn-primary w-full" type="submit" disabled={submission.pending}>
                    {submission.pending ? "Creating account..." : "Create account"}
                  </button>
                </div>
              </fieldset>
            </form>
          <div class="auth-switch">
            <span>Already a user?</span>
            <button class="btn btn-xs btn-link" onClick={props.onGoToLogin}>
              Log In
            </button>
          </div>
      </div>
    </div>
  )
}

export default SignUp
