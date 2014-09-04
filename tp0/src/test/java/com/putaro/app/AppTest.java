package com.putaro.app;

import junit.framework.Test;
import junit.framework.TestCase;
import junit.framework.TestSuite;
import com.putaro.app.Queue;
/**
 * Unit test for simple App.
 */
public class AppTest extends TestCase
{
	final static int  CANT_ITEMS_A_INSERTAR_TEST= Queue.MAX_QUEUE_SIZE;
	
	/**
     * Create the test case
     *
     * @param testName name of the test case
     */
    public AppTest( String testName )
    {
        super( testName );
    }

    /**
     * @return the suite of tests being tested
     */
    public static Test suite()
    { 	
        return new TestSuite( AppTest.class );
    }

    /**
     * Rigourous Test 
     */
    public void testApp()
    {
    	Queue queue = new Queue();
      	testAddSizeTopEnInserciones(queue);
      	testRemoveSizeTopEnBajas(queue);
    }
    
    /**
     * Se chequean los métodos "add", "top" y "size" bajo una serie de inserciones.
     */
    public void testAddSizeTopEnInserciones(Queue queue)
    {
      	for(int i = 0 ; i < CANT_ITEMS_A_INSERTAR_TEST ; i++)
    	{
      		assertEquals(queue.size(),i);
      		try 
      		{
	    		assertEquals(((Integer)queue.top()).intValue() , 0);
    		} catch (Exception e)
    		{
    			assertEquals(e.getClass() , NullPointerException.class);
    		} finally 
    		{
    			queue.add(new Integer(i));
    		}
    	}
    }
    
    /**
     * Se chequean los métodos "remove", "top" y "size" bajo una serie de bajas.
     */
    public void testRemoveSizeTopEnBajas(Queue queue)
    {
    	for(int topCorrecto = 0 ; topCorrecto < CANT_ITEMS_A_INSERTAR_TEST ; topCorrecto++)
    	{
    		assertEquals(queue.size() , CANT_ITEMS_A_INSERTAR_TEST - topCorrecto);
    		try 
    		{    				    		
	    		assertEquals(((Integer)queue.top()).intValue() , topCorrecto);
    		} catch (Exception e)
    		{
    			assertEquals(e.getClass() , NullPointerException.class);
    		}
    		try 
    		{    				    		
	    		queue.remove();
    		} catch (Exception e)
    		{
    			assertEquals(e.getClass() , NullPointerException.class);
    		}
    	}	
    	assertEquals(queue.size() , 0);
    }
    
}
