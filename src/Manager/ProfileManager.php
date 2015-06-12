<?php


namespace Persona\Manager;


use Persona\Repository\GitRepository;

class ProfileManager
{

    private $branch;
    private $repository;
    private $profileName;

    public function getRepository($repository, $profileName, $branch = 'master')
    {
        $this->branch = $branch;
        $this->repository = $repository;
        $this->profileName = $profileName;
        return $this->pattern($repository);
    }


    private function pattern($segment)
    {
        switch (true) {
            case strpos($segment, 'git'):
                return (new GitRepository($this->repository, $this->profileName))->setBranch($this->branch);
                break;
            default:
                return null;
                break;
        }
    }
}