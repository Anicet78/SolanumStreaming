import Password from "./Password"
import Username from "./Username"

type LoginProps = {
  onGoToRegister: () => void;
};

const Login = (props: LoginProps) => {
  return (
    <div class="flex align-items justify-center h-screen">
      <div class="aura aura-dual m-auto">
        <div class="card bg-base-100">
          <div class="card-body">
            <div class="flex justify-center">
              <span class="text-base">Welcome to <b>Solanum Streaming</b></span>
            </div>

            <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
              <legend class="fieldset-legend text-primary">Login</legend>

              <Username/>
              <Password/>

              <button class="btn btn-soft btn-primary mt-4 mb-4" type="submit">Login</button>
              <div class="flex items-center justify-center gap-1">
                <p class="text-xs">Need an account?</p>
                <button class="btn btn-xs btn-link" onClick={props.onGoToRegister}>
                  Sign Up
                </button>
              </div>
            </fieldset>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Login