package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"time"
	"to-do-list/internal/controller"
	"to-do-list/internal/storage"
)

type Controller interface {
	Create(name, description, deadline string) error
	List(orderByAscending bool) error
	Delete(name string) error
}

var (
	toDeleteTaskName        string
	taskName                string
	taskDescription         string
	taskDeadline            string
	sortByDeadlineAscending bool
	ctrl                    Controller
	tenDaysInHours          = 24 * 10 * time.Hour
	defaultDeadlineValue    = time.Now().Add(tenDaysInHours).Format("2006-01-02")
	rootCmd                 = cobra.Command{
		Use:     "root",
		Version: "v1.0",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			const name = "db/storage.json"

			file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
			if errors.Is(err, os.ErrNotExist) {
				err = os.MkdirAll(filepath.Dir(name), os.ModePerm)
				if err != nil {
					log.Println(err)
				}

				file, err = os.Create(name)
				if err != nil {
					log.Fatalln(err)
				}
			} else if err != nil {
				return err
			}

			ctrl = controller.NewController(storage.NewStorage(file))

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("start")
			return nil
		},
	}

	readCommand = &cobra.Command{
		Use:   "read",
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("list")
			orderByAsc, err := cmd.Flags().GetBool("asc")
			if err != nil {
				return err
			}
			err = ctrl.List(orderByAsc)
			if err != nil {
				return err
			}

			return nil
		},
	}

	createCommand = &cobra.Command{
		Use:   "create",
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("create")
			if ctrl == nil {
				return errors.New("nil controller")
			}
			name, err := rootCmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			description, err := rootCmd.Flags().GetString("description")
			if err != nil {
				return err
			}

			deadline, err := rootCmd.Flags().GetString("deadline")
			if err != nil {
				return err
			}

			return ctrl.Create(name, description, deadline)
		},
	}

	deleteCommand = &cobra.Command{
		Use:   "delete",
		Long:  "",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("delete")
			if ctrl == nil {
				return errors.New("nil controller")
			}
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			err = ctrl.Delete(name)
			if err != nil {
				return err
			}

			log.Println("Successfully deleted task")
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(createCommand)
	rootCmd.AddCommand(readCommand)
	rootCmd.AddCommand(deleteCommand)

	readCommand.Flags().BoolVarP(&sortByDeadlineAscending, "asc", "a", false, "Показать сначала старые задачи")
	deleteCommand.PersistentFlags().StringVarP(&toDeleteTaskName, "name", "n", "", "Название задачи, которую надо удалить")

	rootCmd.PersistentFlags().StringVarP(&taskName, "name", "n", "", "Название задачи")
	rootCmd.PersistentFlags().StringVarP(&taskDescription, "description", "d", "", "Описание задачи")
	rootCmd.PersistentFlags().StringVar(&taskDeadline, "deadline", defaultDeadlineValue, "Время окончания задачи")

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln("error ", err)
	}

	log.Println(taskName, taskDescription, taskDeadline)
}
