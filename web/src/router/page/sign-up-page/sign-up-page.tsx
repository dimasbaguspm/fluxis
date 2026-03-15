import { cx } from "@/lib/cx";
import { DEEP_LINKS } from "@constants/page-routes";
import { useRegister } from "@hooks/use-api";
import { useSessionHandler } from "@providers/session";
import { vAlignCenter, vFlex, vFlexCol, vGap2, vGap8, vTextCenter } from "@versaur/core/utilities";
import { AppLayout, ButtonGroup, FormGroup } from "@versaur/react/blocks";
import { EmailInput, PasswordInput, TextInput } from "@versaur/react/forms";
import { Banner, Button, Heading, Text } from "@versaur/react/primitive";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router";

interface SignUpFormInputs {
  displayName: string;
  email: string;
  password: string;
  confirmPassword: string;
}

export const SignUpPage = () => {
  const navigate = useNavigate();
  const { setTokens } = useSessionHandler();
  const [registerUser, err, { isPending }] = useRegister({
    onSuccess: (data) => {
      setTokens(data.accessToken || "", data.refreshToken || "");
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

  const onSubmit = async (data: SignUpFormInputs) => {
    await registerUser({
      displayName: data.displayName,
      email: data.email,
      password: data.password,
    });
  };

  return (
    <AppLayout>
      <AppLayout.Body centered>
        <AppLayout.Main className={cx(vFlex, vFlexCol, vAlignCenter)}>
          <div
            className={cx(vFlex, vFlexCol, vGap8)}
            style={{ maxWidth: "28rem", paddingTop: "3rem" }}
          >
            <div className={cx(vFlex, vFlexCol, vGap2, vTextCenter)}>
              <Heading as="h1" size="2xl">
                Create Account
              </Heading>
              <Text>Join Fluxis today</Text>
            </div>

            <FormGroup onSubmit={handleSubmit(onSubmit)}>
              {err && (
                <FormGroup.Field>
                  <Banner variant="warning">{err.message}</Banner>
                </FormGroup.Field>
              )}

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
                  <Button type="submit" loading={isPending}>
                    {isPending ? "Creating account..." : "Sign Up"}
                  </Button>
                </ButtonGroup>
              </FormGroup.Field>
              <FormGroup.Field>
                <Text as="span">
                  Already have an account? <Link to={DEEP_LINKS.SIGN_IN}>Sign In</Link>
                </Text>
              </FormGroup.Field>
            </FormGroup>
          </div>
        </AppLayout.Main>
      </AppLayout.Body>
    </AppLayout>
  );
};
