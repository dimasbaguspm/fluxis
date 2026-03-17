import { useListProjects } from "@/hooks/use-api";
import { Icon, Loader, Text } from "@versaur/react/primitive";
import { ComboboxInput } from "@versaur/react/forms";
import { useCallback, useMemo, useState } from "react";
import { debounce } from "radash";
import { SearchIcon } from "@versaur/icons";
import { Controller } from "react-hook-form";
import type { Control, FieldValues, Path } from "react-hook-form";

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

  const [projects, err, { isPending }] = useListProjects(query, {
    enabled: searchTerm.length > 0,
  });

  const [selectedProject, setSelectedProject] = useState<any>(null);

  const handleSelectionChange = useCallback(
    (newValue: string | string[] | null, onChange: (value: any) => void) => {
      if (multiple) {
        const values = Array.isArray(newValue) ? newValue : newValue ? [newValue] : [];
        onChange(values);
      } else {
        const value = Array.isArray(newValue) ? newValue[0] : newValue;
        onChange(value || null);
        if (value && projects?.items) {
          const selected = projects.items.find((project) => project.id === value);
          if (selected) {
            setSelectedProject(selected);
          }
        }
      }
    },
    [multiple, projects],
  );

  const getDisplayOptions = useCallback(() => {
    const searchResults = projects?.items || [];
    if (selectedProject && !searchResults.find((project) => project.id === selectedProject.id)) {
      return [selectedProject, ...searchResults];
    }
    return searchResults;
  }, [projects, selectedProject]);

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
              {inputValue.length === 0 ? (
                <Text as="span">Start typing to search projects</Text>
              ) : isPending ? (
                <Loader />
              ) : err ? (
                <Text as="span">Error loading projects</Text>
              ) : getDisplayOptions().length === 0 ? (
                <Text as="span">No projects found</Text>
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
