<?php


namespace Persona\Repository;


use Symfony\Component\Process\ProcessBuilder;

class GitRepository implements RepositoryInterface
{

    private $branch;
    private $repository;
    private $process;
    private $profileName;

    public function __construct($repository, $profileName)
    {
        $this->repository = $repository;
        $this->profileName = $profileName;
    }

    public function setBranch($branch)
    {
        $this->branch = $branch;

        return $this;
    }

    public function getOutput()
    {
        $this->processBuilder();
        $this->process->run(function ($type, $buffer) {
            echo 'git :: '.$buffer;
        });
    }

    public function processBuilder()
    {
        $kernel = new \AppKernel();
        $settings = $kernel->get('settings');
        $builder = new ProcessBuilder();
        $this->process = $builder->setArguments(
            [
                'git',
                'clone',
                '-b',
                $this->branch,
                '--single-branch',
                $this->repository,
                $settings['path_profile'].DIRECTORY_SEPARATOR.$this->profileName
            ])
            ->getProcess();
    }
}