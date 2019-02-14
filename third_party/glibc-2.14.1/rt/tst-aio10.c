#include <aio.h>
#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <unistd.h>

static pthread_barrier_t b;
static pthread_t main_thread;
static int flag;


static void *
tf (void *arg)
{
  int e = pthread_barrier_wait (&b);
  if (e != 0 && e != PTHREAD_BARRIER_SERIAL_THREAD)
    {
      puts ("child: barrier_wait failed");
      exit (1);
    }

  /* There is unfortunately no other way to try to make sure the other
     thread reached the aio_suspend call.  This test could fail on
     highly loaded machines.  */
  sleep (2);

  pthread_kill (main_thread, SIGUSR1);

  while (1)
    sleep (1000);

  return NULL;
}


static void
sh (int sig)
{
  flag = 1;
}


static int
do_test (void)
{
  main_thread = pthread_self ();

  struct sigaction sa;

  sa.sa_handler = sh;
  sa.sa_flags = 0;
  sigemptyset (&sa.sa_mask);

  if (sigaction (SIGUSR1, &sa, NULL) != 0)
    {
      puts ("sigaction failed");
      return 1;
    }

  if (pthread_barrier_init (&b, NULL, 2) != 0)
    {
      puts ("barrier_init");
      return 1;
    }

  int fds[2];
  if (pipe (fds) != 0)
    {
      puts ("pipe failed");
      return 1;
    }

  char buf[42];
  struct aiocb req;
  req.aio_fildes = fds[0];
  req.aio_lio_opcode = LIO_READ;
  req.aio_reqprio = 0;
  req.aio_offset = 0;
  req.aio_buf = buf;
  req.aio_nbytes = sizeof (buf);
  req.aio_sigevent.sigev_notify = SIGEV_NONE;

  pthread_t th;
  if (pthread_create (&th, NULL, tf, NULL) != 0)
    {
      puts ("create failed");
      return 1;
    }

  int e = pthread_barrier_wait (&b);
  if (e != 0 && e != PTHREAD_BARRIER_SERIAL_THREAD)
    {
      puts ("parent: barrier_wait failed");
      exit (1);
    }

  struct aiocb *list[1];
  list[0] = &req;

  e = lio_listio (LIO_WAIT, list, 1, NULL);
  if (e != -1)
    {
      puts ("lio_listio succeeded");
      return 1;
    }
  if (errno != EINTR)
    {
      printf ("lio_listio did not return EINTR: %d (%d = %m)\n", e, errno);
      return 1;
    }

  return 0;
}

#define TEST_FUNCTION do_test ()
#define TIMEOUT 5
#include "../test-skeleton.c"
