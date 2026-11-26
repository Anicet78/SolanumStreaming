import { createSignal, type Component } from "solid-js";
import Password from "./Password"
import PasswordConfirm from "./PasswordConfirm"
import Username from "./Username"
import AuthPanel from "./AuthPanel";
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
      submitLabel="Create account"
      pendingLabel="Creating account..."
      pending={submission.pending}
      error={submission.error ? (submission.error.message ?? "Something went wrong") : undefined}
      onSubmit={handleSubmit}
      switchText="Already a user?"
      switchLabel="Log In"
      onSwitch={props.onGoToLogin}
    >
      <Username value={username()} onInput={setUsername} />
      <Password value={password()} onInput={setPassword} autocomplete="new-password" />
      <PasswordConfirm password={password}/>
    </AuthPanel>
  )
}

export default SignUp
