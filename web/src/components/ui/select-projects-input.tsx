import { useListProjects, useGetProject } from "@/hooks/use-api";
import { Icon, Loader, Text } from "@versaur/react/primitive";
import { ComboboxInput } from "@versaur/react/forms";
import { useCallback, useMemo, useState } from "react";
import { debounce } from "radash";
import { SearchIcon } from "@versaur/icons";
import { Controller, useWatch } from "react-hook-form";
import type { Control, FieldValues, Path } from "react-hook-form";
import { vM2 } from "@versaur/core/utilities";

interface SelectProjectsInputProps<T extends FieldValues> {
  control: Control<T>;
  name: Path<T>;
  label?: string;
  placeholder?: string;
  helper?: string;
  required?: boolean;
  disabled?: boolean;
  multiple?: boolean;
  orgId?: string;
}

export const SelectProjectsInput = <T extends FieldValues>({
  control,
  name,
  label = "Project",
  placeholder = "Search projects...",
  helper,
  required,
  disabled,
  multiple = false,
  orgId,
}: SelectProjectsInputProps<T>) => {
  const [inputValue, setInputValue] = useState("");
  const [searchTerm, setSearchTerm] = useState("");

  const val = useWatch({ control, name });

  const debouncedSearch = useMemo(
    () =>
      debounce({ delay: 300 }, (term: string) => {
        setSearchTerm(term);
      }),
    [],
  );

  const handleSearchChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const value = e.currentTarget.value;
      setInputValue(value);
      debouncedSearch(value);
    },
    [debouncedSearch],
  );

  const query = useMemo(
    () => ({ pageSize: 15, name: searchTerm, ...(orgId ? { orgId: [orgId] } : {}) }),
    [searchTerm, orgId],
  );

  const [projects, err, { isPending }] = useListProjects(query);
  const [project] = useGetProject(val);

  const selectedProject = useMemo(() => {
    if (!val) return null;

    if (projects?.items) {
      const found = projects.items.find((p) => p.id === val);
      if (found) return found;
    }

    return project;
  }, [val, projects?.items, project]);

  const handleSelectionChange = useCallback(
    (newValue: string | string[] | null, onChange: (value: any) => void) => {
      if (multiple) {
        const values = Array.isArray(newValue) ? newValue : newValue ? [newValue] : [];
        onChange(values);
      } else {
        const value = Array.isArray(newValue) ? newValue[0] : newValue;
        onChange(value || null);
      }
    },
    [multiple],
  );

  const getDisplayOptions = useCallback(() => {
    return projects?.items || [];
  }, [projects]);

  return (
    <Controller
      name={name}
      control={control}
      rules={{ required: required ? `${label} is required` : undefined }}
      render={({ field, fieldState: { error } }) => {
        const listboxContent = (
          <>
            <ComboboxInput.Button>
              {multiple
                ? "Select projects"
                : selectedProject
                  ? selectedProject.name
                  : "Select project"}
            </ComboboxInput.Button>
            <ComboboxInput.Container
              variant="list"
              search={
                <ComboboxInput.Search
                  name="project-search"
                  value={inputValue}
                  onChange={handleSearchChange}
                  placeholder="Search projects..."
                />
              }
            >
              {isPending ? (
                <Loader type="bar" />
              ) : err ? (
                <Text as="span" className={vM2}>
                  Error loading projects
                </Text>
              ) : getDisplayOptions().length === 0 ? (
                <Text as="span" className={vM2}>
                  No projects found
                </Text>
              ) : (
                getDisplayOptions().map((project) => (
                  <ComboboxInput.Option key={project.id} value={project.id}>
                    {project.name}
                  </ComboboxInput.Option>
                ))
              )}
            </ComboboxInput.Container>
            {multiple && <ComboboxInput.SelectionChips />}
          </>
        );

        return multiple ? (
          <ComboboxInput
            multiple
            value={(field.value as string[]) || []}
            onChange={(newValue) => handleSelectionChange(newValue, field.onChange)}
            label={label}
            placeholder={placeholder}
            helper={helper}
            required={required}
            disabled={disabled}
            invalid={!!error}
            error={error?.message}
            iconLeft={<Icon as={SearchIcon} />}
          >
            {listboxContent}
          </ComboboxInput>
        ) : (
          <ComboboxInput
            multiple={false}
            value={(field.value as string) || ""}
            onChange={(newValue) => handleSelectionChange(newValue, field.onChange)}
            label={label}
            placeholder={placeholder}
            helper={helper}
            required={required}
            disabled={disabled}
            invalid={!!error}
            error={error?.message}
            iconLeft={<Icon as={SearchIcon} />}
          >
            {listboxContent}
          </ComboboxInput>
        );
      }}
    />
  );
};
