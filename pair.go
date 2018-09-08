// Package pair provides primitives to manage co-author declarations in Git commit templates.
package pair

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gonzalo-bulnes/pair/git"
)

const version = "0.1.0" // adheres to semantic versioning

// Stop removes the co-author declaration from the commit template, if any.
func Stop(out, errors io.Writer) error {
	commitTemplatePath, _, err := git.GetCommitTemplatePath()
	if err != nil {
		switch err.(type) {
		case *git.NoCommitTemplateConfigurationError:
			break
		default:
			return err
		}
	}
	if commitTemplatePath == "" {
		return nil
	}

	config, err := git.NewConfig(commitTemplatePath)
	if err != nil {
		return err
	}

	if author, present := config.CommitTemplate.CoAuthor(); present {
		ok := config.CommitTemplate.RemoveCoAuthor(author)
		if !ok {
			return fmt.Errorf("Unable to remove co-author '%s'", author)
		}
	}
	err = config.Apply()
	if err != nil {
		return err
	}

	return nil
}

// Version prints the package version.
func Version(out, errors io.Writer) error {
	fmt.Fprintf(out, fmt.Sprintf("pair version %s\n", version))
	return nil
}

// With adds a co-author declaration to the current commit template if any,
// or creates a new commit template and configures Git to use it.
func With(out, errors io.Writer, pair string) error {

	err := Stop(out, errors)
	if err != nil {
		return err
	}

	commitTemplatePath, _, err := git.GetCommitTemplatePath()
	if err != nil {
		if err != nil {
			switch err.(type) {
			case *git.NoCommitTemplateConfigurationError:
				break
			default:
				return err
			}
		}
	}
	if commitTemplatePath == "" {
		commitTemplatePath = defaultCommitTemplatePath()
		ensureExists(commitTemplatePath)
		err = git.SetCommitTemplate(commitTemplatePath)
		if err != nil {
			return err
		}
	}
	config, err := git.NewConfig(commitTemplatePath)
	if err != nil {
		return err
	}

	config.CommitTemplate.AddCoAuthor(pair)
	err = config.Apply()
	if err != nil {
		return err
	}
	return nil
}

func defaultCommitTemplatePath() string {
	return filepath.Join(os.Getenv("HOME"), ".pair")
}

func ensureExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			return err
		}
	}
	return nil
}
