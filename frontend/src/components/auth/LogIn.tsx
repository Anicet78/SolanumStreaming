import Password from "./Password"
import Username from "./Username"

type LoginProps = {
  onGoToSignup: () => void;
};

const LogIn = (props: LoginProps) => {
  return (
    <div class="flex align-items justify-center h-screen">
      <div class="aura aura-dual m-auto">
        <div class="card bg-base-100">
          <div class="card-body">
            <div class="flex justify-center">
              <span class="text-base">Welcome to <b>Solanum Streaming</b></span>
            </div>

            <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
              <legend class="fieldset-legend text-primary">Log In</legend>

              <Username/>
              <Password/>

              <button class="btn btn-soft btn-primary mt-4 mb-4" type="submit">Log In</button>
            </fieldset>
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

export default LogIn