<?php


namespace Persona\Command;


use Persona\Command;
use Persona\Manager\ProfileManager;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Filesystem\Filesystem;

class AddProfileCommand extends Command
{
    protected function configure()
    {
        $this->setName('add')
            ->addArgument('profileName', InputArgument::REQUIRED, 'name of profile')
            ->addArgument('repository', InputArgument::REQUIRED, 'remote repository profile')
            ->addOption('branch', 'b', InputOption::VALUE_OPTIONAL, 'branch based', 'master')
            ->setDescription('to add a profile for a remote (git for now)')
        ;
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $profileName = $input->getArgument('profileName');
        $repository = $input->getArgument('repository');
        $branch = $input->getOption('branch');

        $fullPath = $this->get('profiles_directory').DIRECTORY_SEPARATOR.$profileName;
        $fs = new Filesystem();
        if ($fs->exists($fullPath)) {
            $fs->remove($fullPath);
        }
        $profileManager = new ProfileManager();
        $process = $profileManager->getRepository($repository, $profileName, $branch);
        $output->writeln($process->getOutput());

        $settings = $this->get('settings');

        $settings['profiles_repository'][$profileName] = [
            'branch' => $branch,
            'repository' => $repository
        ];

        $fs->dumpFile($this->get('file_settings'), json_encode($settings));
    }

}