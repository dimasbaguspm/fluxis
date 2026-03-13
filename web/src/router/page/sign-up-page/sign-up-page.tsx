import { DEEP_LINKS } from "@constants/page-routes";
import { useRegister } from "@hooks/api/use-auth";
import { useSessionStore } from "@providers/session";
import { AppLayout, ButtonGroup, FormGroup } from "@versaur/react/blocks";
import { EmailInput, PasswordInput, TextInput } from "@versaur/react/forms";
import { Button } from "@versaur/react/primitive";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";

interface SignUpFormInputs {
  displayName: string;
  email: string;
  password: string;
  confirmPassword: string;
}

export const SignUpPage = () => {
  const navigate = useNavigate();
  const setSession = useSessionStore((state) => state.setSession);
  const [, , { isPending }, { mutate }] = useRegister({
    onSuccess: (data) => {
      setSession({
        user: null,
        accessToken: data.accessToken || null,
        refreshToken: data.refreshToken || null,
      });
      navigate(DEEP_LINKS.DASHBOARD);
    },
  });

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<SignUpFormInputs>({
    defaultValues: {
      displayName: "",
      email: "",
      password: "",
      confirmPassword: "",
    },
  });

  const password = watch("password");

  const onSubmit = (data: SignUpFormInputs) => {
    mutate({
      displayName: data.displayName,
      email: data.email,
      password: data.password,
    });
  };

  return (
    <AppLayout>
      <AppLayout.Body centered>
        <AppLayout.Main>
          <h1>Create Account</h1>
          <FormGroup onSubmit={handleSubmit(onSubmit)}>
            <FormGroup.Field>
              <TextInput
                placeholder="Full Name"
                label="Full Name"
                required
                disabled={isPending}
                error={errors.displayName?.message}
                {...register("displayName", {
                  required: "Full name is required",
                })}
              />
            </FormGroup.Field>
            <FormGroup.Field>
              <EmailInput
                placeholder="Email"
                label="Email"
                required
                disabled={isPending}
                error={errors.email?.message}
                {...register("email", {
                  required: "Email is required",
                  pattern: {
                    value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                    message: "Invalid email address",
                  },
                })}
              />
            </FormGroup.Field>
            <FormGroup.Field>
              <PasswordInput
                placeholder="Password"
                label="Password"
                required
                disabled={isPending}
                error={errors.password?.message}
                {...register("password", {
                  required: "Password is required",
                  minLength: {
                    value: 6,
                    message: "Password must be at least 6 characters",
                  },
                })}
              />
            </FormGroup.Field>
            <FormGroup.Field>
              <PasswordInput
                placeholder="Confirm Password"
                label="Confirm Password"
                required
                disabled={isPending}
                error={errors.confirmPassword?.message}
                {...register("confirmPassword", {
                  required: "Confirm password is required",
                  validate: (value) => value === password || "Passwords do not match",
                })}
              />
            </FormGroup.Field>
            <FormGroup.Field>
              <ButtonGroup>
                <Button type="submit" disabled={isPending}>
                  {isPending ? "Creating account..." : "Sign Up"}
                </Button>
              </ButtonGroup>
            </FormGroup.Field>
            <FormGroup.Field>
              <span>
                Already have an account? <a href={DEEP_LINKS.SIGN_IN}>Sign In</a>
              </span>
            </FormGroup.Field>
          </FormGroup>
        </AppLayout.Main>
      </AppLayout.Body>
    </AppLayout>
  );
};
