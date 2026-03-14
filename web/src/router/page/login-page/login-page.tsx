import { DEEP_LINKS } from "@/constants";
import { useSessionHandler } from "@/providers/session";
import { useLogin } from "@hooks/use-api";
import { AppLayout, ButtonGroup, FormGroup } from "@versaur/react/blocks";
import { EmailInput, PasswordInput } from "@versaur/react/forms";
import { Button } from "@versaur/react/primitive";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";

interface LoginFormInputs {
  email: string;
  password: string;
}

export const LoginPage = () => {
  const navigate = useNavigate();
  const { setTokens } = useSessionHandler();
  const [signIn, , { isPending }] = useLogin({
    onSuccess: (data) => {
      setTokens(data.accessToken || null, data.refreshToken || null);
      navigate(DEEP_LINKS.DASHBOARD);
    },
  });

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormInputs>({
    defaultValues: { email: "", password: "" },
  });

  const onSubmit = async (data: LoginFormInputs) => {
    await signIn(data);
  };

  return (
    <AppLayout>
      <AppLayout.Body centered>
        <AppLayout.Main>
          <h1>Sign In</h1>
          <FormGroup onSubmit={handleSubmit(onSubmit)}>
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
                })}
              />
            </FormGroup.Field>
            <FormGroup.Field>
              <ButtonGroup>
                <Button type="submit" disabled={isPending}>
                  {isPending ? "Signing in..." : "Sign In"}
                </Button>
              </ButtonGroup>
            </FormGroup.Field>
          </FormGroup>
        </AppLayout.Main>
      </AppLayout.Body>
    </AppLayout>
  );
};
