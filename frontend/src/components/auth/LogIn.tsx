import { createSignal, type Component } from "solid-js";
import Password from "./Password"
import Username from "./Username"
import AuthPanel from "./AuthPanel";
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
    <AuthPanel
      eyebrow="Welcome back"
      tagline="Your next favorite movie is waiting."
      submitLabel="Log In"
      pendingLabel="Logging in..."
      pending={submission.pending}
      error={submission.error ? (submission.error.message ?? "Something went wrong") : undefined}
      onSubmit={handleSubmit}
      switchText="Need an account?"
      switchLabel="Sign Up"
      onSwitch={props.onGoToSignup}
    >
      <Username value={username()} onInput={setUsername} />
      <Password value={password()} onInput={setPassword} />
    </AuthPanel>
  )
}

export default LogIn;
