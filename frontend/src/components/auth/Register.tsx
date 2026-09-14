import Password from "./Password"
import PasswordConfirm from "./PasswordConfirm"
import Username from "./Username"

type RegisterProps = {
  onGoToLogin: () => void;
};

const Register = (props: RegisterProps) => {
  return (
    <div class="flex align-items justify-center h-screen">
      <div class="aura aura-dual m-auto">
        <div class="card bg-base-100">
          <div class="card-body">
            <div class="flex justify-center">
              <span class="text-base">Welcome to <b>Solanum Streaming</b></span>
            </div>

            <fieldset class="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
              <legend class="fieldset-legend text-primary">Register</legend>

              <Username/>
              <br/>
              <Password/>
              <PasswordConfirm/>

              <button class="btn btn-soft btn-primary mt-4" type="submit">Create account</button>
              <button onClick={props.onGoToLogin}>Back to login</button>
            </fieldset>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Register