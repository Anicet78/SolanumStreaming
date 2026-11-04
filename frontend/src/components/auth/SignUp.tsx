import { createSignal, type Component } from "solid-js";
import Password from "./Password"
import PasswordConfirm from "./PasswordConfirm"
import Username from "./Username"
import AuthPanel from "../AuthPanel";
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
    <AuthPanel
      eyebrow="Start your library"
      tagline="Create an account and keep your watchlist close."
      switchText="Already a user?"
      switchLabel="Log In"
      onSwitch={props.onGoToLogin}
    >
      <form class="grid gap-3" onSubmit={handleSubmit}>
        <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
          <legend class="fieldset-legend text-brand">Sign Up</legend>

          <Username value={username()} onInput={setUsername} />
          <br/>
          <Password value={password()} onInput={setPassword} />
          <PasswordConfirm password={password}/>

          <div class="tooltip tooltip-open tooltip-top tooltip-end mt-4 mb-2" data-tip={submission.error ? submission.error.message : ""}>
            <button class="btn btn-soft btn-primary w-full border-0 bg-brand text-ink" type="submit" disabled={submission.pending}>
              {submission.pending ? "Creating account..." : "Create account"}
            </button>
          </div>
        </fieldset>
      </form>
    </AuthPanel>
  )
}

export default SignUp
