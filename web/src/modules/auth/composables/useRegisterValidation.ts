import useVuelidate from "@vuelidate/core";
import { required, minLength, email, helpers } from "@vuelidate/validators";
import { IRegisterForm } from "@/modules/auth/interfaces";

export function useRegisterValidation(registerForm: IRegisterForm) {
  const rules = {
    name: {
      required: helpers.withMessage("The name is required", required),
    },
    username: {
      required: helpers.withMessage("The username is required", required),
      minLength: helpers.withMessage(
        "The username must have at least 3 characters",
        minLength(3),
      ),
    },
    email: {
      required: helpers.withMessage("The email is required", required),
      email: helpers.withMessage("The email is invalid", email),
    },
    password: {
      required: helpers.withMessage("The password is required", required),
      minLength: helpers.withMessage(
        "The password must have at least 8 characters",
        minLength(8),
      ),
    },
  };

  const v$ = useVuelidate(rules, registerForm);

  return { v$ };
}
