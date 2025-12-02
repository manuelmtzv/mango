import useVuelidate from "@vuelidate/core";
import { required, helpers } from "@vuelidate/validators";
import { ILoginForm } from "@/modules/auth/interfaces";

export function useLoginValidation(loginForm: ILoginForm) {
  const rules = {
    identifier: {
      required: helpers.withMessage("The identifier is required", required),
    },
    password: {
      required: helpers.withMessage("The password is required", required),
    },
  };

  const v$ = useVuelidate(rules, loginForm);

  return { v$ };
}
